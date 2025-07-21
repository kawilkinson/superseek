<?php

namespace App\Http\Middleware;

use Closure;
use Exception;
use Illuminate\Http\Request;
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

    
}