import React from 'react';
import Button from '@material-ui/core/Button';
import Aircraft from '../Models/Aircraft';
import { Checklist } from '../Models/Checklist';
import AddAircraftDialog, { NewAircraft } from './AddAircraftDialog';
import AircraftDetail from './AircraftDetail';
import ChecklistDetail from './ChecklistDetail';

interface MainContainerState {
    aircraft: Aircraft[],
    selectedAircraft: Aircraft | null
    activeChecklist: Checklist | null
    showAddAircraftModal: boolean
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
            activeChecklist: null,
            showAddAircraftModal: false
        }
    }

    // TODO: abstract host address
    componentDidMount() {
        fetch('http://localhost:8080/api/aircraft').then(res => {
            return res.json();
        }).then((data) => {
            this.setState({
                aircraft: data
            });
        }).catch((er: any) => {
            console.error(er);
            this.setState({
                aircraft: []
            });
        });
    }

    toggleSelectedAircraft(id: number) {
        if (this.state.selectedAircraft === null || this.state.selectedAircraft.id !== id) {
            const newSelection = this.state.aircraft.filter(x => x.id === id);
            this.setState({
                selectedAircraft: newSelection[0],
                activeChecklist: null
            });
        }
    }

    // TODO: abstract host address
    setActiveChecklist(id: number) {
        fetch('http://localhost:8080/api/checklists/' + id).then(res => {
            return res.json();
        }).then((response) => {
            this.setState({
                activeChecklist: response
            });
        });
    }

    openChecklistEditor(id: number) {
        if (id !== null) {
            console.log('Open checklist editor for id %s', id);
        }
    }

    openAddAircraftDialog() {
        console.log('Opening the AddAircraft dialog');
        this.setState({
            showAddAircraftModal: true
        });
    }

    onClose() {
        this.setState({
            showAddAircraftModal: false
        });
    }

    async onNewAircraftSubmit(newAC: NewAircraft) {
        let body: any = newAC;
        body.id = 1000;
        await fetch('http://localhost:8080/api/aircraft', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(body),
        }).then(res => res.json())
        .then((response) => {
            console.log(response);
            this.setState({
                showAddAircraftModal: false
            });
        }).catch((e: any) => console.log(e));
    }

    render() {
        return (
            <div id='main-container' style={containerStyle}>
                <div id='selection-div'>
                    {this.state.aircraft.map((ac: Aircraft) => (
                        <AircraftDetail
                            onClick={() => this.toggleSelectedAircraft(ac.id)}
                            key={ac.id}
                            aircraft={ac}
                            selected={this.state.selectedAircraft?.id === ac.id}
                            onClSelected={(cid: number) => this.setActiveChecklist(cid)}
                        />
                    ))}
                    <div id='add-aircraft-btn'>
                        <Button type='button' onClick={() => this.openAddAircraftDialog()}>Add Aircraft</Button>
                    </div>
                </div>
                <div id='active-container'>
                    <ChecklistDetail checklist={this.state.activeChecklist} onEditClick={(chkId: number) => this.openChecklistEditor(chkId) } />
                    {this.state.selectedAircraft ? <div id='add-checklist-btn'>
                        <Button type='button'>Add Checklist</Button>
                    </div> : null }
                </div>
                <AddAircraftDialog open={this.state.showAddAircraftModal} onClose={() => this.onClose()} onSubmit={(a: any) => this.onNewAircraftSubmit(a)} />
            </div>
        )
    }
}