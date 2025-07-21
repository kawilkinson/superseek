<?php

namespace App\Http\Middleware;

use Closure;
use Exception;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Symfony\Component\HttpFoundation\Response;

class FuzzySearch
{
    public function handle(Request $request, Closure $next): Response
    {
        Log::info('FuzzySearch middleware called');
        $query = $request->query('q');
        if (!$query) {
            return $next($request);
        }

        $query = str_replace('+', ' ', $query);
        $queryWords = explode(' ', strtolower($query));
        $processedQuery = [];
        $hasSuggestions = false;

        try {
            foreach($queryWords as $word) {
                if (!trim($word)) {
                    continue;
                }
                $suggestion = $this->checkOrSuggestWord($word);
                Log::info('Processing word: ' . $word . ' => ' . ($suggestion ?? 'no suggestion'));

                if ($suggestion && $suggestion !== $word) {
                    $processedQuery[] = $suggestion;
                    $hasSuggestion = true;
                } else {
                    $processedQuery[] = $word;
                }
            }

            $processedQueryString = implode(' ', $processedQuery);
            $request->merge(['processedQuery' => $processedQueryString]);
            if ($hasSuggestions) {
                $request->merge(['hasSuggestions' => true]);
            }

            Log::info('Processed query: ' . $processedQueryString);
        } catch (Exception $e) {
            Log::error('Spell check error: ' . $e->getMessage());
        }

        return $next($request);
    }

    private function checkOrSuggestWord(string $word): ?string
    {
        $cacheKey = 'spellcheck:' . $word;

        if (Cache::has($cacheKey)) {
            Log::info('Cache hit for word: ' . $word);
            return Cache::get($cacheKey);
        }

        try {
            $collection = DB::connection('mongodb')->table('dictionary');
            $exists = $collection->find($word);
            if ($exists) {
                Cache::put($cacheKey, $word, 3600); // time for cache set to 1 hour with 3600
                Log::info('Word found in DB: ' . $word);
                return $word;
            }

            Log::info('Word not found in DB: ' . $word);

            $length = strlen($word);
            $searchLength = $length - 3 > 0 ? $length - 2 : 1;
            $firstTwoChars = substr($word, 0, $searchLength);

            $cursor = DB::connection('mongodb')
                ->table('dictionary')
                ->raw(function ($collection) use ($firstTwoChars, $length) {
                    return $collection->aggregate([
                        [
                            '$match' => [
                                '_id' => ['$regex' => '^' . $firstTwoChars, '$options' => 'i']
                            ]
                        ],
                        [
                            '$addFields' => [
                                'length' => ['$strLenCP' => '$_id']
                            ]
                        ],
                        [
                            '$match' => [
                                'length' => ['$gte' => $length - 1, '$lte' => $length + 1]
                            ]
                        ]
                    ]);
                });
            
                // Levenshtein distance algorithm
                $bestMatch = null;
                $minDistance = PHP_INT_MAX;
                $wordLength = strlen($word);

                foreach ($cursor as $document) {
                    $candidate = $document ->_id;
                    $candidateLength = strlen($candidate);

                    if (abs($candidateLength - $wordLength) > 2) {
                        continue;
                    }

                    $distance = levenshtein($word, $candidate);

                    $maxDistance = $wordLength <= 4 ? 1 : min(2, floor($wordLength / 4));
                    if ($distance <= $maxDistance && $distance < $minDistance) {
                        $minDistance = $distance;
                        $bestMatch = $candidate;
                    }
                }

                $finalSuggestion = $bestMatch ?? $word;

                Log::info('Best match found ' . $finalSuggestion);

                Cache::put($cacheKey, $finalSuggestion, 3600);

                return $finalSuggestion;
        } catch (Exception $e) {
            Log::error('Suggestion lookup error: ' . $e->getMessage());
            return $word;
        }
    }
}