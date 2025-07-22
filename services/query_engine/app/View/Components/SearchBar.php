<?php

namespace App\View\Components;

use Closure;
use Illuminate\Contracts\View\View;
use Illuminate\View\Component;

class SearchBar extends Component
{
    public function __construct()
    {

    }

    public function render(): View|Closure|string
    {
        return view('components.search-bar');
    }
}