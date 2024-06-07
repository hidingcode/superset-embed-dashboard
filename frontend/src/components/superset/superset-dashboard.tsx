"use client"

import { Component } from "react";
import { embedDashboard } from "@superset-ui/embedded-sdk";
import axios from "axios";
import "./superset-dashboard.css";

function fetchGuestTokenFromBackend(): Promise<string> {
    return new Promise<string>((resolve) => {
        axios.get("http://127.0.0.1:5000/api/guest_token").then((response) => {
            resolve(response.data.token);
        }); 
    })
}

export class SupersetDashboard extends Component {
    state = {
        isLoaded: false
    }

    componentDidMount() {
        this.setState({ isLoaded: true });
        
        embedDashboard({
            id: "c0e94d84-82e6-4e8b-ba23-3e54987094cd", // given by the Superset embedding UI
            supersetDomain: "http://127.0.0.1:8088",
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