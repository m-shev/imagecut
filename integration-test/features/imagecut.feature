# file: features/imagecut.feature

Feature: Imagecut service
  I need to be able to create some calendar event and get notify
  And I need to able to get created event and update it

  Scenario: I send request at first time, the image must be downloaded from origin
  upon repeated request, the picture must be obtained from cache
    When I send request to: "http://imagecut:3000/crop/500/500/?origin=http://image-storage/cover.jpg"
    Then There should be status code 200 and header "X-IMAGECUT-FROM-CACHE" equal "false"
    When I send request to: "http://imagecut:3000/crop/500/500/?origin=http://image-storage/cover.jpg"
    Then There should be status code 200 and header "X-IMAGECUT-FROM-CACHE" equal "true"

  Scenario: I send bad request
      When I send request to: "http://imagecut:3000/crop/x/500/?origin=http://image-storage/cover.jpg"
      Then There should be status code 400

  Scenario: Image-storage does not exist
    When I send request to: "http://imagecut:3000/crop/500/500/?origin=http://some/cover.jpg"
    Then There should be response that contains: "no such host"

  Scenario: Image not found in image-storage
    When I send request to: "http://imagecut:3000/crop/500/500/?origin=http://image-storage/some.jpg"
    Then There should be status code 404

  Scenario: The file is in the wrong format
    When I send request to: "http://imagecut:3000/crop/500/500/?origin=http://image-storage/data.json"
    Then There should be response that contains: "unsupported image format"
