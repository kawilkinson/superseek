<?php

use Illuminate\Support\Facades\Route;

Route::get('/', function () {
    if (config('app.env') === 'local') {
        return response()->json(['message' => 'Welcome to Superseek!']);
    } else {
        return redirect('https://superseek.app');
    }
});
