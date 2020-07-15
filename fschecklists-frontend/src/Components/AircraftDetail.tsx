import React from 'react';

interface AircraftDetailProps {
    aircraft: any
    selected: boolean
    onClick: any
    onClSelected: any
}

export default class AircraftDetail extends React.Component<AircraftDetailProps, {}> {
    render() {
        const checklistTitles = this.props.aircraft.checklists.map((cl: any) => {
            return <li onClick={() => this.props.onClSelected(cl.id)} key={cl.id}>{cl.title}</li>
        });
        return (
            <div id={this.props.aircraft.id} onClick={this.props.onClick}>
                {this.props.selected ?
                <div>
                    <h2>{this.props.aircraft.name}</h2>
                    <ul>
                        {checklistTitles}
                    </ul>
                </div>
                : <h2>{this.props.aircraft.alias}</h2>}
            </div>
        )
    }
}