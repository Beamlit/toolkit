import logging

import pytest

# from beamlit.common.settings import init

logger = logging.getLogger(__name__)

logger.info("Loading conftest.py")

@pytest.fixture
def anyio_backend():
    return 'asyncio'

@pytest.fixture
def settings():
    """Provide initialized settings for tests"""
    print("Initializing settings")
    return {}
    # return init()

@pytest.fixture
def mock_client(mocker):
    """Provide a mocked API client"""
    return mocker.patch('beamlit.authentication.new_client')