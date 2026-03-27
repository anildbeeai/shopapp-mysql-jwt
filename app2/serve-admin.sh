#!/bin/bash
PORT=${1:-4000}
echo "⚙️  Admin Panel running at http://localhost:$PORT"
cd "$(dirname "$0")/admin"
python3 -m http.server $PORT
