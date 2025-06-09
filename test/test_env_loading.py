#!/usr/bin/env python3
"""
Test script to verify .env file loading functionality
"""

import requests
import json
import sys
import time

# Configuration
BASE_URL = "http://localhost:8083"  # Port from .env file
API_KEY = "test-api-key-123"  # API key from .env file

def test_health_check():
    """Test the health check endpoint"""
    print("1. Testing health check endpoint...")
    try:
        response = requests.get(f"{BASE_URL}/health")
        print(f"   Status: {response.status_code}")
        print(f"   Response: {response.json()}")
        assert response.status_code == 200
        print("   ✓ Health check passed\n")
    except Exception as e:
        print(f"   ✗ Health check failed: {e}\n")
        return False
    return True

def test_root_endpoint():
    """Test the root endpoint"""
    print("2. Testing root endpoint...")
    try:
        response = requests.get(f"{BASE_URL}/")
        print(f"   Status: {response.status_code}")
        print(f"   Response: {response.json()}")
        assert response.status_code == 200
        print("   ✓ Root endpoint passed\n")
    except Exception as e:
        print(f"   ✗ Root endpoint failed: {e}\n")
        return False
    return True

def test_models_without_auth():
    """Test models endpoint without authentication"""
    print("3. Testing models endpoint without auth...")
    try:
        response = requests.get(f"{BASE_URL}/v1/models")
        print(f"   Status: {response.status_code}")
        print(f"   Response: {response.json()}")
        assert response.status_code == 401
        print("   ✓ Correctly rejected without auth\n")
    except Exception as e:
        print(f"   ✗ Test failed: {e}\n")
        return False
    return True

def test_models_with_auth():
    """Test models endpoint with authentication"""
    print("4. Testing models endpoint with auth...")
    headers = {
        "Authorization": f"Bearer {API_KEY}"
    }
    try:
        response = requests.get(f"{BASE_URL}/v1/models", headers=headers)
        print(f"   Status: {response.status_code}")
        print(f"   Response: {response.json()}")
        assert response.status_code == 200
        print("   ✓ Models endpoint passed with auth\n")
    except Exception as e:
        print(f"   ✗ Test failed: {e}\n")
        return False
    return True

def test_chat_completions():
    """Test chat completions endpoint"""
    print("5. Testing chat completions endpoint...")
    headers = {
        "Authorization": f"Bearer {API_KEY}",
        "Content-Type": "application/json"
    }
    data = {
        "model": "claude-code",
        "messages": [
            {"role": "user", "content": "Say 'Hello from .env test!'"}
        ],
        "max_tokens": 100
    }
    
    try:
        response = requests.post(
            f"{BASE_URL}/v1/chat/completions",
            headers=headers,
            json=data
        )
        print(f"   Status: {response.status_code}")
        if response.status_code == 200:
            result = response.json()
            print(f"   Model used: {result.get('model', 'N/A')}")
            if 'choices' in result and len(result['choices']) > 0:
                print(f"   Response: {result['choices'][0]['message']['content']}")
            print("   ✓ Chat completions passed\n")
        else:
            print(f"   Response: {response.text}")
            return False
    except Exception as e:
        print(f"   ✗ Test failed: {e}\n")
        return False
    return True

def check_server_logs():
    """Check server logs for .env loading"""
    print("6. Checking server logs for .env loading...")
    try:
        with open('server.log', 'r') as f:
            logs = f.read()
            if 'Listening on' in logs and ':8083' in logs:
                print("   ✓ Server is using port 8083 from .env file")
            else:
                print("   ✗ Server might not be using .env configuration")
            
            if 'API Key authentication enabled' in logs:
                print("   ✓ API key authentication is enabled")
            else:
                print("   ✗ API key might not be loaded from .env")
    except Exception as e:
        print(f"   ⚠ Could not read server logs: {e}")

def main():
    print("=== Testing Claude Code API with .env Loading ===\n")
    print("Expected configuration from .env:")
    print(f"  PORT: 8083")
    print(f"  API_KEY: {API_KEY}")
    print(f"  LOG_LEVEL: debug\n")
    
    # Wait a bit for server to fully start
    print("Waiting for server to start...")
    time.sleep(2)
    
    # Run tests
    tests = [
        test_health_check,
        test_root_endpoint,
        test_models_without_auth,
        test_models_with_auth,
        test_chat_completions,
        check_server_logs
    ]
    
    passed = 0
    for test in tests:
        if test():
            passed += 1
    
    print(f"\n=== Summary: {passed}/{len(tests)} tests passed ===")
    
    if passed == len(tests):
        print("\n✓ All tests passed! The .env file is being loaded correctly.")
        return 0
    else:
        print(f"\n✗ {len(tests) - passed} tests failed.")
        return 1

if __name__ == "__main__":
    sys.exit(main())