from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:4981/openai/v1",
    api_key=""
)

response = client.chat.completions.create(
        model="gemini-advanced",
        messages=[
            {"role": "system", "content": "You are a helpful assistant."},
            {"role": "user", "content": "Hello!"}
        ]
    )

print(response.choices[0].message.content)