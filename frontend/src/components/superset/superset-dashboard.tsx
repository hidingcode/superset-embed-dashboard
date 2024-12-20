"use client"

import { Component } from "react";
import { embedDashboard } from "@superset-ui/embedded-sdk";
import axios from "axios";
import "./superset-dashboard.css";

type SupersetDashboardProps = {
    id: string; // The id provided by the embed configuration UI in Superset
};

const supersetDomain = "http://127.0.0.1:8088";
const backendEndpoint = "http://127.0.0.1:5000/api/guest_token";

function fetchGuestTokenFromBackend(): Promise<string> {
    return new Promise<string>((resolve) => {
        axios.get(backendEndpoint).then((response) => {
            resolve(response.data.token);
        }); 
    })
}

export class SupersetDashboard extends Component<SupersetDashboardProps> {
    state = {
        isLoaded: false
    }

    componentDidMount() {
        this.setState({ isLoaded: true });
        
        embedDashboard({
            id: this.props.id, // given by the Superset embedding UI
            supersetDomain: supersetDomain,
            mountPoint: document.getElementById("superset-dashboard")!, // any html element that can contain an iframe
            fetchGuestToken: () => fetchGuestTokenFromBackend(),
            dashboardUiConfig: {
                hideTitle: true,
                hideChartControls: true,
                filters: {
                  visible: false,
                  expanded: false,
                }
            },
        });
    }

    render() {
        return <div id="superset-dashboard"></div>
    }
}