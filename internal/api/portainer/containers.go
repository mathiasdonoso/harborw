package portainer

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type Hostconfig struct {
	Networkmode string `json:"NetworkMode"`
}

type Labels struct {
	ComDockerStackNamespace   string `json:"com.docker.stack.namespace"`
	ComDockerSwarmNodeId      string `json:"com.docker.swarm.node.id"`
	ComDockerSwarmServiceId   string `json:"com.docker.swarm.service.id"`
	ComDockerSwarmServiceName string `json:"com.docker.swarm.service.name"`
	ComDockerSwarmTask        string `json:"com.docker.swarm.task"`
	ComDockerSwarmTaskId      string `json:"com.docker.swarm.task.id"`
	ComDockerSwarmTaskName    string `json:"com.docker.swarm.task.name"`
	DevopsService             string `json:"devops-service"`
	Maintainer                string `json:"maintainer"`
	Version                   string `json:"version"`
}

type Ipamconfig struct {
	Ipv4address string `json:"IPv4Address"`
}

type Ingress struct {
	Aliases             any        `json:"Aliases"`
	Driveropts          any        `json:"DriverOpts"`
	Endpointid          string     `json:"EndpointID"`
	Gateway             string     `json:"Gateway"`
	Globalipv6address   string     `json:"GlobalIPv6Address"`
	Globalipv6prefixlen int        `json:"GlobalIPv6PrefixLen"`
	Ipamconfig          Ipamconfig `json:"IPAMConfig"`
	Ipaddress           string     `json:"IPAddress"`
	Ipprefixlen         int        `json:"IPPrefixLen"`
	Ipv6gateway         string     `json:"IPv6Gateway"`
	Links               any        `json:"Links"`
	Macaddress          string     `json:"MacAddress"`
	Networkid           string     `json:"NetworkID"`
}

// type Ipamconfig struct {
// 	Ipv4address string `json:"IPv4Address"`
// }

type OmniExecutiveOmniNet struct {
	Aliases             any        `json:"Aliases"`
	Driveropts          any        `json:"DriverOpts"`
	Endpointid          string     `json:"EndpointID"`
	Gateway             string     `json:"Gateway"`
	Globalipv6address   string     `json:"GlobalIPv6Address"`
	Globalipv6prefixlen int        `json:"GlobalIPv6PrefixLen"`
	Ipamconfig          Ipamconfig `json:"IPAMConfig"`
	Ipaddress           string     `json:"IPAddress"`
	Ipprefixlen         int        `json:"IPPrefixLen"`
	Ipv6gateway         string     `json:"IPv6Gateway"`
	Links               any        `json:"Links"`
	Macaddress          string     `json:"MacAddress"`
	Networkid           string     `json:"NetworkID"`
}

type Networks struct {
	Ingress              Ingress              `json:"ingress"`
	OmniExecutiveOmniNet OmniExecutiveOmniNet `json:"omni-executive_omni-net"`
}

type Networksettings struct {
	Networks Networks `json:"Networks"`
}

type ContainersResult struct {
	Command         string          `json:"Command"`
	Created         int             `json:"Created"`
	Hostconfig      Hostconfig      `json:"HostConfig"`
	Id              string          `json:"Id"`
	Image           string          `json:"Image"`
	Imageid         string          `json:"ImageID"`
	Labels          Labels          `json:"Labels"`
	Mounts          []interface{}   `json:"Mounts"`
	Names           []interface{}   `json:"Names"`
	Networksettings Networksettings `json:"NetworkSettings"`
	Ports           []interface{}   `json:"Ports"`
	State           string          `json:"State"`
	Status          string          `json:"Status"`
}

func (p *portainerApiClient) GetContainersJson(endpoint int) (*[]ContainersResult, error) {
	slog.Debug(fmt.Sprintf("Fetching container json from endpoint %d", endpoint))
	url := fmt.Sprintf("%s/api/endpoints/%s/docker/containers/json", p.baseUrl, strconv.Itoa(endpoint))

	fmt.Printf("url is: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.token))
	req.Header.Set("User-Agent", "harborw/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var endpointsResp []ContainersResult
	if err := json.NewDecoder(resp.Body).Decode(&endpointsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	slog.Debug("Container json info fetched", fmt.Sprintf("%+v", endpointsResp))

	return &endpointsResp, nil
}
