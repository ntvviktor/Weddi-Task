import asyncio
import os
import re
import ssl
import time
from bs4 import BeautifulSoup
import certifi
from playwright.async_api import async_playwright
import uuid
import aiohttp
import aiofiles


class ImageScraper:
    def __init__(self, images_dir="../images"):
        self.images_dir = images_dir

    async def save_image(self, review_id: str, filename: str, url: str, ext="png") -> None:
        folder_name = os.path.join(self.images_dir, review_id)
        if not os.path.exists(folder_name):
            os.makedirs(folder_name)
        ssl_context = ssl.create_default_context(cafile=certifi.where())

        async with aiohttp.ClientSession(trust_env=True) as session:
            async with session.get(url, ssl=ssl_context) as response:
                if response.status == 200:
                    content = await response.read()
                    async with aiofiles.open(os.path.join(folder_name, f"{filename}.{ext}"), "wb") as f:
                        await f.write(content)

    async def scrape_review_image(self, review_id: str, review_url: str) -> bool:
        async with async_playwright() as pw:
            browser = await pw.chromium.launch(headless=True)
            context = await browser.new_context(viewport={"width": 1920, "height": 1080})

            start_time = time.time()
            try:
                print("Starting to scrape...")
                page = await context.new_page()
                await page.goto(review_url, timeout=0)

                await page.wait_for_selector(".KtCyie", timeout=0)

                image_section = await page.query_selector(".KtCyie")
                if image_section:
                    expand_button = await page.query_selector(".Tap5If")
                    if expand_button:
                        await page.click('div.Tap5If')

                    images_container = await page.query_selector(".KtCyie")
                    inner_html = await images_container.inner_html()
                    soup = BeautifulSoup(inner_html, 'html.parser')

                    buttons = soup.find_all('button', class_='Tya61d')

                    tasks = [self.process_button(button, review_id) for button in buttons]
                    await asyncio.gather(*tasks)

                await page.close()
                print(f"The time to scrape one vendor: {time.time() - start_time}")

            except Exception as e:
                print(f"An error occurred during the scraping process: {e}")
                return False
            finally:
                await browser.close()

        return True

    async def process_button(self, button, review_id):
        style = button.get('style')
        if style:
            url_start = style.find("url(")
            image_url = style[url_start:]
            pattern = r'url\("([^"]+)"\)'

            match = re.search(pattern, image_url)
            if match:
                url = match.group(1)
                await self.save_image(review_id=review_id, filename=str(uuid.uuid4())[:8], url=url)