import pytest

from translate import translate_resource


def test_translate():
    input_ = {
        "resourceTypeName": "KMSKey",
        "KMSKey": {
            "Description": "string",
            "Enabled": "boolean",
            "KeyPolicy": "object-string",
            "KeyUsage": "KeyUsageEnum",
        },
        "KeyUsageEnum": ["ENCRYPT_DECRYPT", "SIGN_VERIFY", "GENERATE_VERIFY_MAC"],
        "typeName": "AWS::KMS::Key",
    }
    expected = {
        "rust": """
use serde::{Serialize, Deserialize};
use serde_json::Value;

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct KMSKey {
    pub description: Option<String>,
    pub enabled: Option<bool>,
    pub key_policy: Option<Value>,
    pub key_usage: Option<KMSKeyKeyUsageEnum>,
}

enum KMSKeyKeyUsageEnum {
    ENCRYPT_DECRYPT,
    SIGN_VERIFY,
    GENERATE_VERIFY_MAC,
}

""",
        "gql": """
type KMSKey implements Resource {
    CcapiId: String!
    CcapiTypeName: String!
    Description: String
    Enabled: Boolean
    KeyPolicy: String
    KeyUsage: KMSKeyKeyUsageEnum
}
input KMSKeyFilter {
    Description: String
    Enabled: Boolean
    KeyPolicy: String
    KeyUsage: KMSKeyKeyUsageEnum
}

enum KMSKeyKeyUsageEnum {
    ENCRYPT_DECRYPT
    SIGN_VERIFY
    GENERATE_VERIFY_MAC
}

""",
    }
    rust, gql, _ = translate_resource(input_)
    print(gql)
    assert rust == expected["rust"]
    assert gql == expected["gql"]
