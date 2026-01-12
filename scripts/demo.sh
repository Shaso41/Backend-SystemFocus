#!/bin/bash

# Demo script for Redis Clone
# This script demonstrates the key features of the server

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║           Redis Clone - Feature Demonstration            ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo ""

# Check if server is running
echo "📡 Checking server connection..."
if ! nc -z localhost 6379 2>/dev/null; then
    echo "❌ Server is not running on port 6379"
    echo "   Please start the server first: ./redis-clone"
    exit 1
fi

echo "✅ Server is running!"
echo ""

# Function to send command and show response
send_command() {
    echo "▶ $1"
    echo "$2" | nc localhost 6379
    echo ""
    sleep 0.5
}

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1️⃣  Testing Basic Commands"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

send_command "PING" "*1\r\n\$4\r\nPING\r\n"
send_command "SET name 'Redis Clone'" "*3\r\n\$3\r\nSET\r\n\$4\r\nname\r\n\$12\r\nRedis Clone\r\n"
send_command "GET name" "*2\r\n\$3\r\nGET\r\n\$4\r\nname\r\n"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "2️⃣  Testing Key Management"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

send_command "SET user:1 'Alice'" "*3\r\n\$3\r\nSET\r\n\$6\r\nuser:1\r\n\$5\r\nAlice\r\n"
send_command "SET user:2 'Bob'" "*3\r\n\$3\r\nSET\r\n\$6\r\nuser:2\r\n\$3\r\nBob\r\n"
send_command "KEYS *" "*2\r\n\$4\r\nKEYS\r\n\$1\r\n*\r\n"
send_command "EXISTS user:1" "*2\r\n\$6\r\nEXISTS\r\n\$6\r\nuser:1\r\n"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "3️⃣  Testing Expiration"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

send_command "SET session 'xyz123' EX 10" "*5\r\n\$3\r\nSET\r\n\$7\r\nsession\r\n\$6\r\nxyz123\r\n\$2\r\nEX\r\n\$2\r\n10\r\n"
send_command "TTL session" "*2\r\n\$3\r\nTTL\r\n\$7\r\nsession\r\n"
send_command "EXPIRE user:1 30" "*3\r\n\$6\r\nEXPIRE\r\n\$6\r\nuser:1\r\n\$2\r\n30\r\n"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "4️⃣  Testing Delete Operations"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

send_command "DEL user:2" "*2\r\n\$3\r\nDEL\r\n\$6\r\nuser:2\r\n"
send_command "EXISTS user:2" "*2\r\n\$6\r\nEXISTS\r\n\$6\r\nuser:2\r\n"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "5️⃣  Server Information"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

send_command "INFO" "*1\r\n\$4\r\nINFO\r\n"

echo "✅ Demo completed successfully!"
echo ""
echo "💡 Tip: You can connect manually using:"
echo "   telnet localhost 6379"
echo "   or"
echo "   redis-cli -p 6379"
