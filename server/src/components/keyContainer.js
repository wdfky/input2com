
import { useEffect, useState } from "react"
import js2hid from "./js-to-hid.json"
import { DndRecipient, ItemTypes } from "./dnd"
import { DndProvider, useDrag, useDrop } from 'react-dnd';



const KeyContainer = ({ className, setKeyboardConfig ,onKeyboardClick, code, children }) => {

    const [{ isOver, canDrop }, drop] = useDrop({
        accept: ItemTypes.CARD,
        drop: (item, monitor) => {
            // console.warn(item, monitor);
            setKeyboardConfig(Number(js2hid[code]).toString() ,item["key"] )
        },
        hover: (item, monitor) => {
            console.log(item, monitor);
        },
        collect: (monitor) => ({
            isOver: !!monitor.isOver(),
            canDrop: !!monitor.canDrop(),
        }),
        canDrop: (item, monitor) => true
    });


    const [state, setState] = useState(false)
    const onClick = (e) => {
        onKeyboardClick(Number(js2hid[code]).toString() , e.clientX , e.clientY)
    }

    return <div
        className={className}
        onClick={onClick}
        ref={drop}
        style={
            state ? {
                backgroundColor: "#ff7e27ff"
            } : {}
        }

    >{children}</div>
}

export default KeyContainer;