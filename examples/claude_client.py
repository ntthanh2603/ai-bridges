from langchain_anthropic import ChatAnthropic

# Initialize the client pointing to our local bridge
llm = ChatAnthropic(
    base_url="http://localhost:4981/claude", 
    model="claude-3-5-sonnet-20240620",
    temperature=0.7,
    api_key="abc"
)
response = llm.invoke("Hello!")
print(response.content)
