#!/bin/bash

# Примеры запросов к TON Wallet API
# Использование: bash examples.sh

BASE_URL="http://localhost:8080"

echo "=========================================="
echo "TON Wallet API - Примеры запросов"
echo "=========================================="
echo ""

# 1. Health Check
echo "1. Health Check"
echo "GET $BASE_URL/health"
curl -X GET "$BASE_URL/health"
echo -e "\n"

# 2. Создать кошелек
echo "=========================================="
echo "2. Создать новый кошелек"
echo "POST $BASE_URL/api/v1/wallet"
curl -X POST "$BASE_URL/api/v1/wallet" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "wallet_type": "V5R1Final",
    "network": "testnet"
  }'
echo -e "\n"

# 3. Получить информацию о кошельке (ID=1)
echo "=========================================="
echo "3. Получить информацию о кошельке (ID=1)"
echo "GET $BASE_URL/api/v1/wallet/1"
curl -X GET "$BASE_URL/api/v1/wallet/1"
echo -e "\n"

# 4. Получить баланс кошелька (ID=1)
echo "=========================================="
echo "4. Получить баланс кошелька (ID=1)"
echo "GET $BASE_URL/api/v1/wallet/1/balance"
curl -X GET "$BASE_URL/api/v1/wallet/1/balance"
echo -e "\n"

# 5. Список кошельков пользователя (user_id=1)
echo "=========================================="
echo "5. Список кошельков пользователя (user_id=1)"
echo "GET $BASE_URL/api/v1/wallet/list?user_id=1"
curl -X GET "$BASE_URL/api/v1/wallet/list?user_id=1"
echo -e "\n"

# 6. Удалить кошелек (ID=1) - раскомментируйте при необходимости
# echo "=========================================="
# echo "6. Удалить кошелек (ID=1)"
# echo "DELETE $BASE_URL/api/v1/wallet/1"
# curl -X DELETE "$BASE_URL/api/v1/wallet/1"
# echo -e "\n"

echo "=========================================="
echo "Тестирование завершено!"
echo "=========================================="
