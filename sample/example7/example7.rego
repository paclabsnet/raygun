package example7

import rego.v1


default allow := false

    # Always allow OPTIONS calls
    allow := true if {
        print("input is: ",input)
        input.attributes.request.http.method == "OPTIONS"
    }

    # Allow GET requests to /get endpoint for everyone
    allow := true if {
        input.attributes.request.http.method == "GET"
        input.attributes.request.http.path == "/get"
    }

    # Allow access to /headers for users with valid API key
    allow := true if {
        input.attributes.request.http.method == "GET"
        input.attributes.request.http.path == "/headers"
        input.attributes.request.http.headers["x-api-key"] == "demo-key-123"
    }

    # Allow access to /headers if the IP address is in an approved VPN subnet
    #
    allow := true if {
        input.attributes.request.http.method == "GET"
        input.attributes.request.http.path == "/headers"
        "x-ipaddress" in object.keys(input.attributes.request.http.headers)
	vpn_user(input.attributes.request.http.headers["x-ipaddress"])
    }

    # Allow POST to /post for admin users
    #
    allow := true if {
        input.attributes.request.http.method == "POST"
        input.attributes.request.http.path == "/post"
        input.attributes.request.http.headers["x-user-role"] == "admin"
    }

    # Deny access to /status/* endpoints
    allow := true if {
        input.attributes.request.http.method == "GET"
        startswith(input.attributes.request.http.path, "/status/")
        false  # This will never be true, effectively denying access
    }

 
    # allow POST messages if the caller is in the approved US states
    #
    allow := true if {
        input.attributes.request.http.method == "POST"
        input.attributes.request.http.path == "/post"
        location := input.attributes.request.http.headers["x-location"]
        # location in [ "california", "colorado", "nevada" ]
	location in data.envoy.regulations.approved_locations
    }
        

#
# is the provided ip address contained within any of the CIDRs 
# associated with VPNs?
#
vpn_user( ip_address ) := true if {
	some subnet in data.envoy.network.subnets.vpn
	net.cidr_contains(subnet, ip_address)
} else := false 
