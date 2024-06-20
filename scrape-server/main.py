import asyncio
from flask import Flask, Response, jsonify, request

from image_scraper import ImageScraper

app = Flask(__name__)

class ScraperAPI():
    # Contructor for ScrapeAPI
    def __init__(self) -> None:
        self.app = Flask(__name__)
        self.register_endpoints()
    
    @property
    def app(self) -> Flask:
        return self._app

    @app.setter
    def app(self, app: Flask) -> None:
        self._app = app

    def register_endpoints(self) -> None:
        self.app.add_url_rule(rule="/v1/scrape", endpoint="scrape", view_func=self.scrape, methods=["POST"])

    def scrape(self) -> Response:
        data = request.get_json()
        review_id = data["review_id"]
        review_url = data["review_url"]
        scraper = ImageScraper()
        success = asyncio.run(scraper.scrape_review_image(review_id, review_url))
        
        if success:
            return jsonify({"status": "ok", "message": "Scraping successful"})
        else:
            return jsonify({"status": "error", "message": "Scraping failed"})

    def run(self, debug=True, port=5000) -> None:
        self.app.run(debug=debug, port=port)

if __name__ == "__main__":
    scraperAPI = ScraperAPI()
    scraperAPI.run()