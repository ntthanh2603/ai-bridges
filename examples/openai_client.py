from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:4981/openai/v1",
    api_key=""
)

response = client.chat.completions.create(
        model="gpt-3.5-turbo",
        messages=[
            {"role": "system", "content": "You are a helpful assistant."},
            {"role": "user", "content": "Hello!"}
        ]
    )

print(response.choices[0].message.content)