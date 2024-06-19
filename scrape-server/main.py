import asyncio
from flask import Flask, jsonify, request

from async_scrape_v2 import scrape_review_image

app = Flask(__name__)

@app.route("/api/scrape", methods=["POST"])
def scrape():
    data = request.get_json()
    review_id = data["review_id"]
    review_url = data["review_url"]
    success = asyncio.run(scrape_review_image(review_id, review_url))
    
    if success:
        return jsonify({"status": "OK", "message": "Scraping successful"}), 200
    else:
        return jsonify({"status": "error", "message": "Scraping failed"}), 500
    

def main():
    app.run(debug=True)


if __name__ == "__main__":
    main()