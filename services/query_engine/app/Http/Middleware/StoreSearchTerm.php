<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use Symfony\Component\HttpFoundation\Response;

class StoreSearchTerm
{
    public function handle(Request $request, Closure $next): Response
    {
        $searchTerm = $request->get('processedQuery');

        Log::info('StoreSearchTerm middleware called.');
        Log::info('Search term stored: ' . ($searchTerm ?? 'null'));

        if (!$searchTerm) {
            return $next($request);
        }
        $searchTerm = trim($searchTerm);

        Redis::zincrby('top_searches', 1, strtolower($searchTerm));
        Redis::zremrangebyrank('top_searches', 0, -101);
        Redis::incr('total_searches');
        Redis::expire('total_searches', 86400);

        return $next($request);
    }
}