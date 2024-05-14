# fantasy-shop-api
This project was built to practice constructing CRUD API endpoints using the clean architecture style. Additionally, it incorporates Google OAuth2 for user authentication and includes some unit testing.

## About this project
The Fantasy Shop API is a RESTful API that enables users to perform CRUD (Create, Read, Update, Delete) operations, with certain routes protected by OAuth2 that only allow admins or authorized users access.

## Overview
* **GET**
    *   /v1/health (check endpoint health)
    *   /v1/item-shop (get all items in shop)
    *   /v1/inventory (get all items in player's inventory)
* **POST**
    *  /v1/item-shop/buy (buy item only players are allowed)
    *  /v1/item-shop/sell (sell item only players are allowed)
    *  /v1/item-managing (create item only admins are allowed)
* and much more..