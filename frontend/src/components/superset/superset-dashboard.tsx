"use client"

import * as React from "react";
import { Component } from "react";
import { embedDashboard } from "@superset-ui/embedded-sdk";
import axios from "axios";
import "./superset-dashboard.css";

interface SupersetDashboardProps {
    id: string; // The id provided by the embed configuration UI in Superset
};

const supersetDomain = "http://127.0.0.1:8088";
const backendEndpoint = "http://127.0.0.1:5000/api/guest_token";

async function fetchGuestTokenFromBackend(): Promise<string> {
    const response = await axios.get<{ token: string }>(backendEndpoint);
    return response.data.token;
}

export class SupersetDashboard extends Component<SupersetDashboardProps> {
    state = {
        isLoaded: false
    }

    componentDidMount(): void {
        this.setState({ isLoaded: true });
        
        void embedDashboard({
            id: this.props.id, // given by the Superset embedding UI
            supersetDomain,
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
        return <div id="superset-dashboard" />;
    }
}