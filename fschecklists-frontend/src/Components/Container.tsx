import React from 'react';
import AircraftDetail from './AircraftDetail';
import Aircraft from '../Models/Aircraft';
import { Checklist } from '../Models/Checklist';
import ChecklistDetail from './ChecklistDetail';

interface MainContainerState {
    aircraft: Aircraft[],
    selectedAircraft: Aircraft | null
    activeChecklist: Checklist | null
}

const containerStyle = {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-around'
} as React.CSSProperties

export default class Container extends React.Component<{}, MainContainerState> {
    constructor(props: any) {
        super(props);
        this.state = {
            aircraft: [],
            selectedAircraft: null,
            activeChecklist: null
        }

        this.toggleSelectedAircraft.bind(this);
    }

    componentDidMount() {
        fetch('http://localhost:8080/api/aircraft').then(res => {
            return res.json();
        }).then((data) => {
            this.setState({
                aircraft: data
            });
        });
    }

    toggleSelectedAircraft(id: number) {
        if (this.state.selectedAircraft === null || this.state.selectedAircraft.id !== id) {
            const newSelection = this.state.aircraft.filter(x => x.id === id);
            this.setState({
                selectedAircraft: newSelection[0]
            });
        } else {
            this.setState({
                selectedAircraft: null
            });
        }
    }

    setActiveChecklist(id: number) {
        fetch('http://localhost:8080/api/checklists/' + id).then(res => {
            return res.json();
        }).then((response) => {
            this.setState({
                activeChecklist: response
            });
        });
    }

    render() {
        return (
            <div id='main-container' style={containerStyle}>
                <div id='selection-div'>
                    {this.state.aircraft.map((ac: any) => (
                        <AircraftDetail
                            onClick={() => this.toggleSelectedAircraft(ac.id)}
                            key={ac.id}
                            aircraft={ac}
                            selected={this.state.selectedAircraft?.id === ac.id}
                            onClSelected={(cid: number) => this.setActiveChecklist(cid)}
                        />
                    ))}
                </div>
                <div id='active-container'>
                    <ChecklistDetail checklist={this.state.activeChecklist} />
                </div>
            </div>
        )
    }
}