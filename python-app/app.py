import os
from flask import Flask, render_template, request, redirect, url_for
import redis

app = Flask(__name__)

# --- Redis Connection ---
# Get connection details from environment variables set in docker-compose.yml
REDIS_HOST = os.environ.get('REDIS_HOST', 'localhost') # Default to localhost if not in Docker
REDIS_PORT = int(os.environ.get('REDIS_PORT', 6379))
REDIS_QUEUE_NAME = os.environ.get('REDIS_QUEUE_NAME', 'garden_actions') # Match the name used by the Go worker

r = redis.StrictRedis(host=REDIS_HOST, port=REDIS_PORT, decode_responses=True) # decode_responses=True for getting strings back

# --- Routes ---

@app.route('/')
def index():
    """Serve the main garden interaction page."""
    return render_template('index.html')

@app.route('/action', methods=['POST'])
def perform_action():
    """Handle button clicks and send action to Redis."""
    print("hello")
    action_type = request.form.get('action') # Get the value from the clicked button
    
    if action_type:
        print(f"Received action: {action_type}. Sending to Redis queue: {REDIS_QUEUE_NAME}")
        try:
            # Push the action message onto the Redis list
            r.rpush(REDIS_QUEUE_NAME, action_type)
            print(f"Action '{action_type}' sent to Redis.")
        except Exception as e:
            print(f"Error sending action to Redis: {e}")
            # In a real app, handle errors gracefully (e.g., show error message)

    # Redirect back to the main page after handling the action
    return redirect(url_for('index'))

# --- Run the app ---
if __name__ == '__main__':
    # Use 0.0.0.0 to make the server accessible from outside the container
    app.run(host='0.0.0.0', port=5000, debug=True) # debug=True is useful during development
