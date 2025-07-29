# Superseek

[![CI](https://github.com/kawilkinson/search-engine/actions/workflows/CI.yaml/badge.svg)](https://github.com/kawilkinson/search-engine/actions/workflows/CI.yaml)

## Introduction
Welcome to SuperSeek, the 16 bit search engine!

## Services and Databases
Superseek uses a microservice architecture and has these services:
1. Crawler
2. Indexer
3. Image Indexer
4. TF-IDF
5. PageRank
6. Query Engine
7. Client

It also uses these two databases to help store and create queues for data:
1. MongoDB
2. Redis

All of this combined allows Superseek to crawl web pages, index them, give them a weight/score based on an old PageRank algorithm that Google used to use, use TF-IDF to further influence the weight/score, and have a query engine for the client to talk to.
