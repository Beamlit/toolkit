import logging

import pytest
from langchain_openai import ChatOpenAI

from beamlit.agents.decorator import agent
from beamlit.common.settings import init
from beamlit.models import Agent, AgentSpec, EnvironmentMetadata

logger = logging.getLogger(__name__)

init()

def test_agent_decorator_simple():
    """Test basic agent decorator functionality"""
    @agent(
        agent={
            "metadata": {"name": "test_agent"},
            "spec": {"description": "Test agent"}
        },
        override_model=ChatOpenAI(model="gpt-4o-mini")
    )
    def sample_function():
        return True

    assert callable(sample_function)

def test_agent_decorator_must_provide_model():
    """Test basic agent decorator functionality"""
    with pytest.raises(ValueError) as exc_info:
        @agent(agent={
            "metadata": {"name": "test_agent"},
            "spec": {"description": "Test agent"}
        })
        def sample_function():
            return True

    assert "You must provide a model" in str(exc_info.value)

def test_agent_decorator_invalid_input():
    """Test agent decorator with invalid input"""

    with pytest.raises(Exception) as exc_info:
        @agent(agent=Agent(
            metadata=EnvironmentMetadata(name="test"),
            spec=AgentSpec(description="test")
        ))
        def sample_function():
            return True

    assert "agent must be a dictionary" in str(exc_info.value)

@pytest.mark.parametrize("model_name", [
    "non_existent_model",
    "",
    None
])
def test_agent_decorator_missing_model(model_name):
    """Test agent decorator with missing or invalid model"""

    with pytest.raises(Exception) as exc_info:
        @agent(agent={
            "metadata": {"name": "test_agent"},
            "spec": {
                "description": "Test agent",
                "model": model_name
            }
        })
        def sample_function():
            return True

@pytest.fixture
def anyio_backend():
    return 'asyncio'

@pytest.mark.parametrize('anyio_backend', [])
async def test_agent_decorator():
    """Test agent decorator with a valid model"""
    @agent(agent={
        "metadata": {"name": "test_agent"},
        "spec": {"description": "Test agent", "model": "sandbox-openai"}
    })
    async def sample_function(request, agent):
        print(agent)

    await sample_function()