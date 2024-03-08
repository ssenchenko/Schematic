#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionCode {
    image_uri: String,
    zip_file: String,
    s3_object_version: String,
    s3_key: String,
    s3_bucket: String,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionDeadLetterConfig {
    target_arn: String,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionEnvironment {
    variables: String,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionEphemeralStorage {
    size: i32,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionFileSystemConfig {
    local_mount_path: String,
    arn: String,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionTracingConfig {
    mode: LambdaFunctionTracingConfigModeEnum,
}

enum LambdaFunctionTracingConfigModeEnum {
    Active(String),
    PassThrough(String),
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionVpcConfig {
    ipv6_allowed_for_dual_stack: bool,
    subnet_ids: Vec<String>,,
    security_group_ids: Vec<String>,,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionTag {
    value: String,
    key: String,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionImageConfig {
    working_directory: String,
    command: Vec<String>,,
    entry_point: Vec<String>,,
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionSnapStart {
    apply_on: LambdaFunctionSnapStartApplyOnEnum,
}

enum LambdaFunctionSnapStartApplyOnEnum {
    PublishedVersions(String),
    None(String),
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionSnapStartResponse {
    optimization_status: LambdaFunctionSnapStartResponseOptimizationStatusEnum,
    apply_on: LambdaFunctionSnapStartResponseApplyOnEnum,
}

enum LambdaFunctionSnapStartResponseOptimizationStatusEnum {
    On(String),
    Off(String),
}

enum LambdaFunctionSnapStartResponseApplyOnEnum {
    PublishedVersions(String),
    None(String),
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionRuntimeManagementConfig {
    runtime_version_arn: String,
    update_runtime_on: LambdaFunctionRuntimeManagementConfigUpdateRuntimeOnEnum,
}

enum LambdaFunctionRuntimeManagementConfigUpdateRuntimeOnEnum {
    Auto(String),
    FunctionUpdate(String),
    Manual(String),
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunctionLoggingConfig {
    system_log_level: LambdaFunctionLoggingConfigSystemLogLevelEnum,
    application_log_level: LambdaFunctionLoggingConfigApplicationLogLevelEnum,
    log_format: LambdaFunctionLoggingConfigLogFormatEnum,
    log_group: String,
}

enum LambdaFunctionLoggingConfigSystemLogLevelEnum {
    DEBUG(String),
    INFO(String),
    WARN(String),
}

enum LambdaFunctionLoggingConfigApplicationLogLevelEnum {
    TRACE(String),
    DEBUG(String),
    INFO(String),
    WARN(String),
    ERROR(String),
    FATAL(String),
}

enum LambdaFunctionLoggingConfigLogFormatEnum {
    Text(String),
    JSON(String),
}

#[derive(Serialize, Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
struct LambdaFunction {
    arn: String,
    code: LambdaFunctionCode,
    dead_letter_config: LambdaFunctionDeadLetterConfig,
    description: String,
    environment: LambdaFunctionEnvironment,
    ephemeral_storage: LambdaFunctionEphemeralStorage,
    file_system_configs: Vec<LambdaFunctionFileSystemConfig>,,
    function_name: String,
    handler: String,
    architectures: Vec<LambdaFunctionArchitecturesEnum>,,
    kms_key_arn: String,
    layers: Vec<String>,,
    memory_size: i32,
    reserved_concurrent_executions: i32,
    role: String,
    runtime: String,
    tags: Vec<LambdaFunctionTag>,,
    timeout: i32,
    tracing_config: LambdaFunctionTracingConfig,
    vpc_config: LambdaFunctionVpcConfig,
    code_signing_config_arn: String,
    image_config: LambdaFunctionImageConfig,
    package_type: LambdaFunctionPackageTypeEnum,
    policy: String,
    snap_start: LambdaFunctionSnapStart,
    snap_start_response: LambdaFunctionSnapStartResponse,
    runtime_management_config: LambdaFunctionRuntimeManagementConfig,
    logging_config: LambdaFunctionLoggingConfig,
}

enum LambdaFunctionArchitecturesEnum {
    x86_64(String),
    arm64(String),
}

enum LambdaFunctionPackageTypeEnum {
    Image(String),
    Zip(String),
}

