using System.Collections;
using UnityEngine;
using UnityEngine.Networking;

public class NetworkManager : MonoBehaviour
{
    public static NetworkManager Instance;

    public string username;
    public string lobbyId;
    private string server = "http://localhost:8000/gnh?e=";

    private void Awake()
    {
        if (Instance != null && Instance != this)
        {
            Destroy(gameObject);
            return;
        }

        Instance = this;
        DontDestroyOnLoad(gameObject);
    }

	void Start()
	{
        Debug.Log("Signed in as: " + username);
	}

	public void Host()
    {
        StartCoroutine(SendRequest("host+" + username, (response) =>
        {
            lobbyId = response;
            Debug.Log("Hosting Lobby: " + lobbyId);
        }));
    }

    public void Join(string joinLobbyId)
    {
        lobbyId = joinLobbyId;

        StartCoroutine(SendRequest("join+" + lobbyId + "+" + username, (response) =>
        {
            Debug.Log("Joined Lobby: " + lobbyId);
        }));
    }

    public void UpdatePlayer(int hp, Vector3 pos)
    {
        StartCoroutine(SendRequest($"psync+{lobbyId}+{username}+{hp}+{pos.x}+{pos.y}+{pos.z}", (response) =>
        {
            Debug.Log("Sent Player Data.");
        }));
    }

    private IEnumerator SendRequest(string query, System.Action<string> callback)
    {
        using UnityWebRequest www = UnityWebRequest.Get(server + UnityWebRequest.EscapeURL(query));
        yield return www.SendWebRequest();

        if (www.result == UnityWebRequest.Result.Success)
        {
            callback?.Invoke(www.downloadHandler.text);
        }
        else
        {
            Debug.LogWarning("Network Error: " + www.error);
        }
    }
}
