# file: features/imagecut.feature

Feature: Imagecut service
    I need to be able to create some calendar event and get notify
    And I need to able to get created event and update it

    Scenario: Crop image and add to cache
        When I send request to imagecut service: "http://imagecut:3000/crop/500/500/?origin=http://image-store/cover.jpg"
        Then response status is "200" and header "X-IMAGECUT-FROM-CACHE" equal "false"

