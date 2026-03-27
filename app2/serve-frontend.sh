#!/bin/bash
PORT=${1:-3000}
echo "🌐 Frontend running at http://localhost:$PORT"
cd "$(dirname "$0")/frontend"
python3 -m http.server $PORT
