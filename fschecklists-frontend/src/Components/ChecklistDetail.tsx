import React from 'react';
import { Checklist } from '../Models/Checklist';
import Button from '@material-ui/core/Button'

interface ChecklistProps {
    checklist: Checklist | null
    onEditClick: any
}

export default class ChecklistDetail extends React.Component<ChecklistProps, {}> {
    render() {
        const title = this.props.checklist ? this.props.checklist.title : '';
        const items = this.props.checklist?.items.map((item: any) => {
            return <li key={item}>{item}</li>
        });
        return (
            <div id="checklist-container">
                <div>{title} {this.props.checklist ? 
                    <span><hr /> <Button type='button' onClick={() => this.props.onEditClick(this.props.checklist?.id)}>Edit Checklist</Button></span>
                    : null}
                </div>
                <ul>
                    {items}
                </ul>
            </div>
        )
    }
}
