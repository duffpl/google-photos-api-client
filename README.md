# google-photos-api-client
Yet Another Google Photos API Client

## Description
I needed Google Photos API client for one of my projects. Since Google [removed auto-generated client for Go](https://code-review.googlesource.com/c/google-api-go-client/+/39951) I've decided to create one. 
There are some [other](https://github.com/nmrshll/google-photos-api-client-go) libraries floating around but wanted to give it a try on my own and I just didn't like using generated client from [mirrored generated library](https://github.com/gphotosuploader/googlemirror).

I've tried to make client less complicated compared to generated/existing ones (hopefully).

It seems that API is pretty stable so I guess that client should be working in unforseen future too.

Besides basic API communication all methods have their async/sync wrappers for consuming results on the go and not caring about pagination. 
 
Uploading media items has basic implementation - no partial/resumable uploads. To make life a bit easier there is method that allows batch uploading files specified by path into choosen album.

Error handling might be a bit flaky but should be OK.

[godoc documentation](https://godoc.org/github.com/duffpl/google-photos-api-client)

List of implemented endpoints

### Albums
Based on https://developers.google.com/photos/library/reference/rest/v1/albums
* [x] [albums.addEnrichment](https://developers.google.com/photos/library/reference/rest/v1/albums/addEnrichment)
* [x] [albums.batchAddMediaItems](https://developers.google.com/photos/library/reference/rest/v1/albums/batchAddMediaItems)
* [x] [albums.batchRemoveMediaItems](https://developers.google.com/photos/library/reference/rest/v1/albums/batchRemoveMediaItems)
* [x] [albums.create](https://developers.google.com/photos/library/reference/rest/v1/albums/create)
* [x] [albums.get](https://developers.google.com/photos/library/reference/rest/v1/albums/get)
* [x] [albums.list](https://developers.google.com/photos/library/reference/rest/v1/albums/list)
* [x] [albums.patch](https://developers.google.com/photos/library/reference/rest/v1/albums/patch)
* [x] [albums.share](https://developers.google.com/photos/library/reference/rest/v1/albums/share)
* [x] [albums.unshare](https://developers.google.com/photos/library/reference/rest/v1/albums/unshare)
### Media items
* [x] [mediaItems.batchCreate](https://developers.google.com/photos/library/reference/rest/v1/mediaItems/batchCreate)
* [x] [mediaItems.batchGet](https://developers.google.com/photos/library/reference/rest/v1/mediaItems/batchGet)
* [x] [mediaItems.get](https://developers.google.com/photos/library/reference/rest/v1/mediaItems/get)
* [x] [mediaItems.list](https://developers.google.com/photos/library/reference/rest/v1/mediaItems/list)
* [x] [mediaItems.patch](https://developers.google.com/photos/library/reference/rest/v1/mediaItems/patch)
* [x] [mediaItems.search](https://developers.google.com/photos/library/reference/rest/v1/mediaItems/search)
### Shared albums
* [x] [sharedAlbums.get](https://developers.google.com/photos/library/reference/rest/v1/sharedAlbums/get)
* [x] [sharedAlbums.join](https://developers.google.com/photos/library/reference/rest/v1/sharedAlbums/join)
* [x] [sharedAlbums.leave](https://developers.google.com/photos/library/reference/rest/v1/sharedAlbums/leave)
* [x] [sharedAlbums.list](https://developers.google.com/photos/library/reference/rest/v1/sharedAlbums/list)

## Usage

Client requires prepared HTTP OAuth client.
Example implementation using `golang.org/x/oauth2/google`:
```go
func main() {
    credsConfig, _ := ioutil.ReadFile("google-api-credentials.json")
    oauthConfig, _ := google.ConfigFromJSON(
        credsConfig,
        "https://www.googleapis.com/auth/photoslibrary",
        "profile",
    )
    apiToken := &oauth2.Token{
        RefreshToken: "<< refresh token >>",
    }
    oauthHttpClient := oauthConf.Client(context.Background(), apiToken)
    apiClient := google_photos_api_client.NewApiClient(oauthHttpClient)
}
```
Resources can be accessed through their services in client. E.g.:
```go
...
apiClient := google_photos_api_client.NewApiClient(oauthHttpClient)
items, err := apiClient.MediaItems.ListAll(nil, ctx)
albums, err := apiClient.Albums.ListAll(
    &google_photos_api_client.AlbumsListOptions{
        PageSize:                 20,
        ExcludeNonAppCreatedData: true,
    }, ctx)
...
```

## To do
- [ ] functional tests that'll check if API didn't change
- [ ] unit tests
- [x] implement all endpoints (almost there)
- [ ] add validation to request options structs 