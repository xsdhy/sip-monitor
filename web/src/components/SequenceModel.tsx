import SequenceDiagram from "../views/details/SequenceDiagram";
import {Button} from "antd";

interface Prop {
    callID: string
}

export function SequenceModel(p: Prop) {
    return (
        <div>
            <SequenceDiagram callID={p.callID}/>
            <Button type="link">
                <a target="_blank" href={'/call/details?sip_call_id=' + p.callID}>新页面中打开</a>
            </Button>
        </div>
    );
}