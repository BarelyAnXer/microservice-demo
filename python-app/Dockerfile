# ./python-app/Dockerfile
# Use an official Python runtime as a parent image
FROM python:3.9-slim 

# Set the working directory in the container
WORKDIR /app

# Copy the requirements file into the working directory
COPY requirements.txt .

# Install any needed packages specified in requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the application code into the working directory
COPY . .

# Make port 5000 available to the outside world (default Flask port)
EXPOSE 5000

# Define environment variables for the app
# These will be overridden by docker-compose environment section
ENV REDIS_HOST="localhost"
ENV REDIS_PORT="6379"
ENV REDIS_QUEUE_NAME="garden_actions" # Default name

# Run app.py when the container launches
CMD ["python", "app.py"]