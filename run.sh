#!/bin/bash
echo "====== RUNNING API ======"
cd api
make dev
cd ..
sleep 30
echo "====== MIGRATE SCRAPER DATA TO MONGODB ======"
cd scraper
make import
cd ..
echo "====== RUNNING TEST API ======"
cd api
make test
