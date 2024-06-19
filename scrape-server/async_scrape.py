import asyncio
import csv
import os
import re
import time
from bs4 import BeautifulSoup
from playwright.async_api import async_playwright
from typing import List
import uuid
import requests


def read_csv(csv_filename: str) -> List[dict]:
    """
    Function to read a csv file into an list of objects (dictionaries)
    - params: csv_filename: string

    """
    with open(csv_filename, "r") as csv_file:
        reader = csv.reader(csv_file)

        result = []
        for i, v in enumerate(reader):
            obj = {}
            # Skip the row that contains column names
            if i == 0:
                continue
            else:
                obj["VendorID"] = v[0]
                obj["LinkToSource"] = v[-1]
                result.append(obj)
        return result


def save_image(folder_name: str, filename:str, url: str, ext="png") -> None:
    if not os.path.exists(folder_name):
        os.makedirs(folder_name)
    
    response = requests.get(url)
    with open(os.path.join(folder_name, f"{filename}.{ext}"), "wb") as f:
        f.write(response.content)


async def main():
    vendor_data = read_csv("data.csv")
    async with async_playwright() as pw:
        # Create an instance of a Chromium browser
        browser = await pw.chromium.launch(headless=True)
        context = await browser.new_context(viewport={"width": 1920, "height": 1080})

        """
        Two cases to consider:
        - Whether there is an image section (Google map stored as list of button)
        - Whether there is an expand button (+k images)
        """
        start_time = time.time()
        try:
            for obj in vendor_data:
                vendor_id = obj["VendorID"]
                vendor_source = obj["LinkToSource"]
                print("Starting to scrape...")
                page = await context.new_page()
                await page.goto(vendor_source, timeout=0)

                await page.wait_for_selector(".KtCyie", timeout=0)
                
                image_section = await page.query_selector(".KtCyie")
                # First case: there is an image section
                if image_section:
                    expand_button = await page.query_selector(".Tap5If")
                    # Second case: there is an expand section
                    if expand_button:
                        await page.click('div.Tap5If')
                    # Doesn't matter if there is a expand button, we gonna crawl all 
                    # images, but the expand button need to be clicked if it appear.
                    images_container = await page.query_selector(".KtCyie")
                    inner_html = await images_container.inner_html()
                    soup = BeautifulSoup(inner_html, 'html.parser')
                
                    # Find all buttons with the class Tya61d
                    buttons = soup.find_all('button', class_='Tya61d')
                

                    # TODO: Optimize this
                    for button in buttons:
                        style = button.get('style')
                        if style:
                            # Extract URL from the style attribute
                            url_start = style.find("url(") 
                            image_url = style[url_start:]
                            # Regex pattern to extract URL
                            pattern = r'url\("([^"]+)"\)'

                            # Find the URL using the regex pattern
                            match = re.search(pattern, image_url)

                            # Extracted URL
                            if match:
                                url = match.group(1)
                                save_image(folder_name=vendor_id, filename=str(uuid.uuid4())[:8], url=url)
                                
                await page.close()
                break

            print(f"The time to scrape one vendor: {time.time() - start_time}")
                
        finally:
            await browser.close()


if __name__ == "__main__":
    asyncio.run(main())
