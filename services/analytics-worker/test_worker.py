import json

import pytest

from worker import parse_click_event


def test_parse_click_event():
    body = json.dumps(
        {
            "code": "abc1234",
            "original_url": "https://example.com",
            "clicked_at": "2026-08-24T21:00:00Z",
        }
    ).encode()

    code, original_url, clicked_at = parse_click_event(body)

    assert code == "abc1234"
    assert original_url == "https://example.com"
    assert clicked_at == "2026-08-24T21:00:00Z"


def test_parse_click_event_missing_field_raises_key_error():
    body = json.dumps({"code": "abc1234"}).encode()

    with pytest.raises(KeyError):
        parse_click_event(body)


def test_parse_click_event_invalid_json_raises():
    with pytest.raises(json.JSONDecodeError):
        parse_click_event(b"not json")
