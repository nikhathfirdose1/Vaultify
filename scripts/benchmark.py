import asyncio
import time
import aiohttp
import random
import string

API_URL = "http://localhost:8080"
TOKEN = "abc123"  # Must match your vaultify.yml

# Generates random key-value secrets
def random_secret():
    name = ''.join(random.choices(string.ascii_lowercase, k=8))
    value = ''.join(random.choices(string.ascii_letters + string.digits, k=32))
    return name, value

async def store_secret(session, name, value):
    payload = {"name": name, "value": value, "ttl": 300}
    headers = {"Authorization": f"Bearer {TOKEN}", "Content-Type": "application/json"}
    async with session.post(f"{API_URL}/store", json=payload, headers=headers) as resp:
        return await resp.text(), resp.status

async def fetch_secret(session, name):
    headers = {"Authorization": f"Bearer {TOKEN}"}
    async with session.get(f"{API_URL}/fetch/{name}", headers=headers) as resp:
        return await resp.text(), resp.status

async def benchmark(n_requests=1000):
    async with aiohttp.ClientSession() as session:
        start = time.perf_counter()
        
        tasks = []
        keys = []

        # Store phase
        for _ in range(n_requests):
            name, value = random_secret()
            keys.append(name)
            tasks.append(store_secret(session, name, value))
        results = await asyncio.gather(*tasks)
        print(f"Stored {len(results)} secrets.")

        # Fetch phase
        fetch_tasks = [fetch_secret(session, k) for k in keys]
        fetch_results = await asyncio.gather(*fetch_tasks)
        print(f"Fetched {len(fetch_results)} secrets.")

        duration = time.perf_counter() - start
        print(f"\nBenchmark complete in {duration:.2f} seconds")
        successes = sum(1 for _, status in fetch_results if status == 200)
        print(f"Successful fetches: {successes} / {n_requests}")

if __name__ == "__main__":
    asyncio.run(benchmark(1000))
