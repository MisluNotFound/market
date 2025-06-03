def get_value(data, keys): 
    for key in keys: 
        value = data.get(key) 
        if value is not None: 
            return value 
    return "" 

import json 

def main(arg: str, node: str) -> dict:
    data = json.loads(arg) 
    detail_str = data.get("detail", "")
    detail = json.loads(detail_str) if detail_str else {}

    return { 
        "alertName": get_value(data, ["alertname", "alertName", "alert_name"]),
        "service": get_value(data, ["svc_name", "service"]), 
        "endpoint": get_value(data, ["endpoint", "content_key"]), 
        "pod": get_value(data, ["pod", "src_pod", "pod_name"]), 
        "namespace": get_value(data, ["namespace", "src_namespace"]),
        "pid": get_value(data, ["pid"]),
        "containerId": get_value(data, ["containerId"]),
        "nodeName": node if node else get_value(data, ["node", "src_node", "node_name", "nodeName"]),
        "sourceFrom": get_value(detail, ["source_from"]),
    }