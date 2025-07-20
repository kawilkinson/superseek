<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Redis;
use App\Http\Controllers\QuerySearchController;
use Illuminate\Support\Facades\Cache;

class RedisController extends Controller
{
    public function getTopSearches(Request $request)
    {
        $topSearches = Redis::zrevrange('top_searches', 0, -1);
        if (!$topSearches) {
            return response()->json(['searches' => []]);
        }

        $topTenSearches = array_slice($topSearches, 0, 10);

        return response()->json([
            'searches' => $topTenSearches
        ]);
    }

    public function getSearchSuggestions(Request $request)
    {
        $searchTerm = $request->get('q');

        $topSearches = Redis::zrevrange('top_searches', 0, -1);
        if (!$topSearches) {
            return response()->json(['searches' => []]);
        }

        $suggestions = array_filter($topSearches, function ($search) use ($searchTerm) {
            return stripos($search, $searchTerm) === 0;
        });

        $suggestions = array_slice($suggestions, 0, 10);

        return response()->json([
            'searches' => array_values($suggestions)
        ]);
    }

    public function returnRandomPage(Request $request)
    {
        $topSearches = Redis::zrevrange('top_searches', 0, -1);

        $totalSearches = Redis::get('total_searches');

        $querySearchController = new QuerySearchController();

        if (Cache::has('random_page')) {
            $randomPage = Cache::get('random_page');
        } else {
            $randomPage = $querySearchController->getRandomPage($request);

            Cache::put('random_page', $randomPage, 1440);
        }

        return view('random-results', [
            'topSearches' => $topSearches,
            'totalSearches' => $totalSearches,
            'randomPage' => $randomPage,
        ]);
    }
}