<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class QuerySearchController extends Controller 
{
    public function get_page_connections(Request $request)
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
}