<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class QuerySearchController extends Controller 
{
    public function getPageConnections(Request $request)
    {
        $url = $request->input('url');
        if (!$url) {
            return response()->json(['error' => 'URL is required'], 400);
        }

        error_log('URL: ' . $url);

        $outlinksData = DB::connection('mongodb')
            ->table('outlinks')
            ->where('id', $url)
            ->first();

        $outlinks = $outlinksData->links ?? [];

        $enrichedOutlinks = [];
        if (count($outlinks) > 0) {
            $metadataCollection = DB::connection('mongodb')
                ->table('metadata')
                ->whereIn('_id', $outlinks)
                ->get();
            
            $metadataMap = [];
            foreach ($metadataCollection as $metadata) {
                $metadataMap[$metadata->id] = $metadata;
            }

            foreach ($outlinks as $link) {
                if (isset($metadataMap[$link])) {
                    $enrichedOutlinks[] = [
                    'url' => $link, 
                    'title' => $metadataMap[$link]->title ?? 'Page Not Indexed'
                    ];
                } else {
                    $enrichedOutlinks[] = [
                        'url' => $link,
                        'title' => 'Page Not Indexed'
                    ];
                }
            }
        }
        
        $backlinksData = DB::connection('mongodb')
            ->table('backlinks')
            ->where('id', $url)
            ->first();
        
        $backlinks = $backlinksData->links ?? [];

        $enrichedBacklinks = [];
        if (count($backlinks) > 0) {
            $metadataCollection = DB::connection('mongodb')
                ->table('metadata')
                ->whereIn('_id', $backlinks)
                ->get();
            
            $metadataMap = [];
            foreach ($metadataCollection as $metadata) {
                $metadataMap[$metadata->id] = $metadata;
            }

            foreach ($backlinks as $link) {
                if (isset($metadataMap[$link])) {
                    $enrichedBacklinks[] = [
                        'url' => $link,
                        'title' => $metadataMap[$link]->title ?? 'Page Not Indexed'
                    ];
                } else {
                    $enrichedBacklinks[] = [
                        'url' => $link,
                        'title' => 'Page Not Indexed'
                    ];
                }
            }
        }

        $urlMetadata = DB::connection('mongodb')
            ->table('metadata')
            ->where('_id', $url)
            ->first();
        
        return view('page-connections', [
            'url' => $url,
            'title' => $urlMetadata->title ?? 'Page Not Indexed',
            'outlinks' => $enrichedOutlinks,
            'backlinks' => $enrichedBacklinks,
        ]);
    }

    public function getTopImages($query, $page = 1, $perPage = 5)
    {
        $query = str_replace('+', ' ', $query);
        $words = explode(' ', strtolower($query));

        $countPipeline = [
            ['$match' => ['word' => ['$in' => $words]]],
            ['$group' => ['_id' => '$url']],
            ['$count' => 'total']
        ];

        $countResult = DB::connection('mongodb')
            ->table('word_images')
            ->raw(fn($collection) => $collection->aggregate($countPipeline)->toArray());

        $totalResults = isset($countResult[0]) ? $countResult[0]['total'] : 0;

        $paginationPipeline = [
            ['$match' => ['word' => ['$in' => $words]]],
            [
                '$group' => [
                    '_id' => '$url',
                    'cumWeight' => ['$sum' => '$weight'],
                    'matchedWords' => ['$addToSet' => '$word'],
                    'matchCount' => ['$sum' => 1]
                ]
            ],
            ['$sort' => ['matchCount' => -1, 'cumWeight' => -1]],
            ['$skip' => ($page - 1) * $perPage],
            ['$limit' => $perPage]
        ];

        /** @var array $paginatedResults */
        $paginatedResults = DB::connection('mongodb')
            ->table('word_images')
            ->raw(function ($collection) use ($paginationPipeline) {
                $cursor = $collection->aggregate($paginationPipeline, ['cursor' => ['batchSize' => 20]]);
                $results = [];
                foreach ($cursor as $document) {
                    $results[] = $document;
                }
                return $results;
            });
        
        $urls = array_map(fn($result) => $result['_id'], $paginatedResults);

        $imagesData = DB::connection('mongodb')->table('images')
            ->whereIn('_id', $urls)
            ->get();

        $imageDataByUrl = [];
        foreach ($imagesData as $data) {
            $imageDataByUrl[$data->id] = $data;
        }

        $pageUrls = [];
        foreach ($imageDataByUrl as $result) {
            $pageUrls[] = $result->page_url ?? '';
        }

        $pageMetadataList = DB::connection('mongodb')->table('metadata')
            ->whereIn('_id', $pageUrls)
            ->get();

        $pageMetadataByUrl = [];
        foreach ($pageMetadataList as $meta) {
            $pageMetadataByUrl[$meta->id] = $meta;
        }

        foreach ($paginatedResults as &$result) {
            $imageData = $imageDataByUrl[$result['_id']] ?? null;
            $result['alt'] = $imageData->alt ?? '';
            $result['filename'] = $imageData->filename ?? '';
            $result['page_url'] = $imageData->page_url ?? '';
            $pageMetadata = $pageMetadataByUrl[$result['page_url']] ?? null;
            $result['page_title'] = $pageMetadata->title ?? '';
            $result['page_text'] = '';
            $length = 100;
            if (isset($pageMetadata->summary_text)) {
                $result['page_text'] = strlen($pageMetadata->summary_text) > $length
                    ? substr($pageMetadata->summary_text, 0, $length) . '...'
                    : $pageMetadata->summary_text;
            }
        }

        return [$paginatedResults, $totalResults];
    }

    public function stats()
    {
        $results = DB::connection('mongodb')->table('metadata')->count();

        return response()->json([
            'status' => 'up',
            'pages' => $results,
        ]);
    }

    public function search(Request $request)
    {
        $hasSuggestions = $request->input('hasSuggestions');
        $originalQuery = $request->input('q');
        $processedQuery = $request->input('processedQuery');
        $query = $processedQuery;
        if (!$query) {
            $query = "";
            return view('search-results', [
                'query' => $query,
                'results' => [],
                'total' => 0,
                'topImages' => [],
                'suggestions' => $hasSuggestions,
                'originalQuery' => $originalQuery,
                'page' => 0,
            ]);
        }

        $query = str_replace('+', ' ', $query);
        $words = explode(' ', strtolower($query));

        $perPage = 20;
        $page = $request->input('page', 1);

        $countPipeline = [
            ['$match' => ['word' => ['$in' => $words]]],
            ['$group' => ['_id' => '$url']],
            ['$count' => 'total']
        ];

        $countResult = DB::connection('mongodb')
            ->table('words')
            ->raw(fn($collection) => $collection->aggregate($countPipeline)->toArray());
        
        $totalResults = isset($countResult[0]) ? $countResult[0]['total'] : 0;

        $paginationPipeline = [
            ['$match'=> ['word' => ['$in' => $words]]],
            [
                '$group' => [
                    '_id' => '$url',
                    'cumWeight' => ['$sum' => '$weight'],
                    'matchedWords' => ['$addToSet' => '$word'],
                    'matchCount' => ['$sum' => 1]
                ]
            ],
            ['$sort' => ['matchCount' => -1, 'cumWeight' => -1]],
            ['$skip' => ($page - 1) * $perPage],
            ['$limit' => $perPage]
        ];

        /** @var array $paginatedResults */
        $paginatedResults = DB::connection('mongodb')
            ->table('words')
            ->raw(function ($collection) use ($paginationPipeline) {
                $cursor = $collection->aggregate($paginationPipeline, ['cursor' => ['batchSize' => 20]]);
                $results = [];
                foreach ($cursor as $document) {
                    $results[] = $document;
                }
                return $results;
            });
        
        $urls = array_map(fn($result) => $result['_id'], $paginatedResults);

        $pageRank = DB::connection('mongodb')->table('pagerank')
            ->whereIn('_id', $urls)
            ->get();
        
        error_log('Page rank: ' . json_encode($pageRank));

        $metadata = DB::connection('mongodb')->table('metadata')
            ->whereIn('_id', $urls)
            ->get();
        
        $metadataByUrl = [];
        foreach ($metadata as $meta) {
            $metadataByUrl[$meta->id] = $meta;
        }

        foreach ($paginatedResults as &$result) {
            $resultMetadata = $metadataByUrl[$result['_id']] ?? null;
            $result['description'] = $resultMetadata->description ?? '';
            $result['last_crawled'] = $resultMetadata->last_cralwed ?? '';
            $result['summary_text'] = $resultMetadata->summary_text ?? '';
            $result['title'] = $resultMetadata->title ?? '';

            $result['pagerank'] = $pageRankByUrl[$result['_id']] ?? 0;

            $tfidfWeight = $result['cumWeight'];
            $pageRankWeight = $result['pagerank'];

            $combinedScore = (0.6 * $tfidfWeight) + (0.4 * $pageRankWeight);

            $result['combinedScore'] = $combinedScore;
        }

        usort($paginatedResults, function ($a, $b) {
            return $b['combinedScore'] <=> $a['combinedScore'];
        });

        $topImages = [];
        if ($page == 1) {
            [$topImages, $unused] = $this->getTopImages($query, $page, 5);
        }

        return view('search-results', [
            'query' => $query,
            'results' => $paginatedResults,
            'total' => $totalResults,
            'topImages' => $topImages,
            'suggestions' => $hasSuggestions,
            'originalQuery' => $originalQuery,
            'page' => $page,
        ]);
    }

}