Module beamlit.api.knowledgebases.delete_knowledgebase
======================================================

Functions
---------

`asyncio(knowledgebase_name: str, *, client: beamlit.client.AuthenticatedClient) ‑> beamlit.models.knowledgebase.Knowledgebase | None`
:   Delete knowledgebase
    
     Deletes an knowledgebase by Name.
    
    Args:
        knowledgebase_name (str):
    
    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.
    
    Returns:
        Knowledgebase

`asyncio_detailed(knowledgebase_name: str, *, client: beamlit.client.AuthenticatedClient) ‑> beamlit.types.Response[beamlit.models.knowledgebase.Knowledgebase]`
:   Delete knowledgebase
    
     Deletes an knowledgebase by Name.
    
    Args:
        knowledgebase_name (str):
    
    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.
    
    Returns:
        Response[Knowledgebase]

`sync(knowledgebase_name: str, *, client: beamlit.client.AuthenticatedClient) ‑> beamlit.models.knowledgebase.Knowledgebase | None`
:   Delete knowledgebase
    
     Deletes an knowledgebase by Name.
    
    Args:
        knowledgebase_name (str):
    
    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.
    
    Returns:
        Knowledgebase

`sync_detailed(knowledgebase_name: str, *, client: beamlit.client.AuthenticatedClient) ‑> beamlit.types.Response[beamlit.models.knowledgebase.Knowledgebase]`
:   Delete knowledgebase
    
     Deletes an knowledgebase by Name.
    
    Args:
        knowledgebase_name (str):
    
    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.
    
    Returns:
        Response[Knowledgebase]