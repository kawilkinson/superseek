<?php

use App\Http\Controllers\QuerySearchController;
use App\Http\Controllers\RedisController;
use App\Http\Middleware\FuzzySearch;
use App\Http\Middleware\StoreSearchTerm;
use Illuminate\Support\Facades\Route;

Route::prefix('api')->controller(QuerySearchController::class)->group(function () {
    Route::get('/search', 'search')->name('search')->middleware([FuzzySearch::class, StoreSearchTerm::class]);
    Route::get('/search-images', 'searchImages')->name('search-images');
    Route::get('/dictionary', 'getDictionary')->name('dictionary');
    Route::get('/stats', 'stats')->name('stats');
    Route::get('/top-ranked-pages', 'getTopRankedPage')->name('top-ranked-pages');
    Route::get('/page-connections', 'getPageConnections')->name('page-connections');
});

Route::prefix('api')->controller(RedisController::class)->group(function () {
    Route::get('/get-top-searches', 'getTopSearches')->name('get-top-searches');
    Route::get('/get-search-suggestions', 'getSearchSuggestions')->name('get-search-suggestions');
    Route::get('/random', 'returnRandomPage')->name('random');
});

Route::get('/secret', function () {
    return response()->json(['message'=> 'You have found a secret message, shh, its a secret, do not tell anyone else']);
});