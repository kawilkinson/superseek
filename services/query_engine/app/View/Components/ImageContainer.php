<?php

namespace App\View\Components;

use Closure;
use Illuminate\Contracts\View\View;
use Illuminate\View\Component;

class ImageContainer extends Component
{
    public $url;
    public $alt;
    public $title;
    public $page_url;
    public $text;

    public function __construct($url, $alt, $title, $page, $text)
    {   
        $this->url = $url;
        $this->alt = $alt;
        $this->title = $title;
        $this->page_url = $page;
        $this->text = $text;
    }

    public function render(): View|Closure|string
    {
        return view('components.image-container');
    }
}