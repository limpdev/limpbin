# Recursive downloading using CLI tools -> ripgrep, xidel, curl
import sys 
import re
import requests
from bs4 import BeautifulSoup
from urllib.parse import urljoin
import os
import subprocess

def get_links(url):
    try:
        response = requests.get(url)
        response.raise_for_status()  # Raise HTTPError for bad resp.
        soup = BeautifulSoup(response.content, 'html.parser')
        links = []
        for a_tag in soup.find_all('a', href=True):
            absolute_url = urljoin(url, a_tag['href'])
            links.append(absolute_url)
        return list(set(links))  # Remove duplicate links
    except requests.exceptions.RequestException as e:
        print(f"Error fetching {url}: {e}")
        return []
    except Exception as e:
        print(f"Error parsing {url}: {e}")
        return []

if __name__ == "__main__":
    base_url = sys.argv[1]
    # Updated regex: capture only the domain
    match = re.search(r'^https?://([^/]+)', base_url)
    if match:
        domain = match.group(1)
        print(f"Domain extracted: {domain}")
    else:
        print("Invalid URL format. Using default filename 'output.txt'.")
        domain = "output"
    
    # Create directory for markdown files if it doesn't exist
    os.makedirs(domain, exist_ok=True)
    
    file_name = domain + ".txt"
    link_list = get_links(base_url)
    
    with open(file_name, "w") as f:
        for link in link_list:
            f.write(f"{link}\n")

    # Now we open each link in the list and download the content as markdown
    moka = "moka.exe"
    for i, link in enumerate(link_list, 1):
        # Create a unique filename for each markdown file
        # Uses link name or falls back to a sequential number
        try:
            # Try to extract a meaningful filename from the URL
            match = re.search(r'/([^/]+)$', link)
            if match:
                filename = match.group(1)
                # Remove any non-alphanumeric characters except periods and hyphens
                filename = re.sub(r'[^\w\-\.]', '', filename)
                # Truncate to a reasonable length
                filename = filename[:50]
            else:
                filename = f"page_{i}"
            
            # Construct full path for markdown file
            linkName = os.path.join(os.getcwd(), domain, f"{filename}.md")
            
            # Ensure unique filename if it already exists
            base, ext = os.path.splitext(linkName)
            counter = 1
            while os.path.exists(linkName):
                linkName = f"{base}_{counter}{ext}"
                counter += 1
            
            # Run Moka with input link and output markdown file
            subprocess.run([moka, link, linkName], check=True)
            print(f"Converted {link} to {linkName}")
        
        except Exception as e:
            print(f"Error converting {link}: {e}")