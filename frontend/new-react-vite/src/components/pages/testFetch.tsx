const BASE_URL = 'http://example.org';

//         Input T ↴   is thread through to ↴
async function api<T>(path: string): Promise<T> {
    const response = await fetch(`${BASE_URL}/${path}`);

    if (!response.ok) {
      throw new Error(response.statusText);
    }

    //    And can also be used here ↴
    return await response.json() as T;
}

// Set up various fetches
async function getConfig() {
  return await api<{ version: number }>('config');
}

// Elsewhere
async function main() {
  const config = await getConfig();
  console.log(config.version); 
}