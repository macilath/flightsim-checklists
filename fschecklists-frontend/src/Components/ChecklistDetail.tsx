import React from 'react';
import { Checklist } from '../Models/Checklist';

interface ChecklistProps {
    checklist: Checklist | null
}

export default class ChecklistDetail extends React.Component<ChecklistProps, {}> {
    render() {
        const title = this.props.checklist ? this.props.checklist.title : '';
        const items = this.props.checklist?.items.map((item: any) => {
            return <li key={item}>{item}</li>
        });
        return (
            <div id="checklist-container">
                <div>{title}</div>
                <ul>
                    {items}
                </ul>
            </div>
        )
    }
}
