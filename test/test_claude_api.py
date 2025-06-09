#!/usr/bin/env python3
"""
Test script for Claude Code API
This script demonstrates how to use the OpenAI client library
to connect to the Claude Code API server.
"""

import openai
import json

# Configure the client to use Claude Code API
client = openai.OpenAI(
    api_key="test-api-key",  # Your API key
    base_url="http://localhost:8082/v1"  # Claude Code API server
)

def test_basic_chat():
    """Test basic chat completion"""
    print("1. Testing basic chat completion...")
    response = client.chat.completions.create(
        model="claude-code",
        messages=[
            {"role": "user", "content": "Say hello and tell me what day it is"}
        ]
    )
    print(f"Response: {response.choices[0].message.content}\n")

def test_file_operations():
    """Test file operations with auto permissions"""
    print("2. Testing file operations with auto permissions...")
    response = client.chat.completions.create(
        model="claude-code",
        messages=[
            {"role": "user", "content": "Create a file called test_output.txt with the content 'Claude Code API works!' in the current directory"}
        ],
        extra_body={
            "claude_options": {
                "tools": ["Write"],
                "working_dir": "/tmp",
                "auto_allow_permissions": True
            }
        }
    )
    print(f"Response: {response.choices[0].message.content}\n")

def test_code_execution():
    """Test code execution with Bash"""
    print("3. Testing code execution...")
    response = client.chat.completions.create(
        model="claude-code",
        messages=[
            {"role": "user", "content": "List all Python files in the current directory"}
        ],
        extra_body={
            "claude_options": {
                "tools": ["Bash"],
                "working_dir": "/tmp",
                "auto_allow_permissions": True
            }
        }
    )
    print(f"Response: {response.choices[0].message.content}\n")

def test_file_analysis():
    """Test file reading and analysis"""
    print("4. Testing file analysis...")
    response = client.chat.completions.create(
        model="claude-code",
        messages=[
            {"role": "user", "content": "Read this Python script and explain what it does"}
        ],
        extra_body={
            "claude_options": {
                "tools": ["Read"],
                "files": ["/tmp/test_claude_api.py"],
                "auto_allow_permissions": True
            }
        }
    )
    print(f"Response: {response.choices[0].message.content}\n")

def main():
    """Run all tests"""
    print("=== Claude Code API Test Suite ===\n")
    
    try:
        test_basic_chat()
        test_file_operations()
        test_code_execution()
        test_file_analysis()
        
        print("=== All tests completed! ===")
        
    except Exception as e:
        print(f"Error: {e}")
        print(f"Make sure the Claude Code API server is running on http://localhost:8082")

if __name__ == "__main__":
    main()