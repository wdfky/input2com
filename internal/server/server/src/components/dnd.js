import { Box, Grid, Paper } from "@mui/material";
import { HTML5Backend } from 'react-dnd-html5-backend';
import { DndProvider, useDrag, useDrop } from 'react-dnd';
import { useEffect } from "react";


const ItemTypes = {
    CARD: 'card',
};

const DndSender = ({ args, type, children }) => {
    const [{ isDragging }, drag] = useDrag({
        type: type,
        item: args,
        collect: (monitor) => ({
            isDragging: !!monitor.isDragging(),
        }),
    });

    return (
        <div
            ref={drag}
            className={type}
            style={{
                opacity: isDragging ? 0.5 : 1,
                cursor: 'move',
                width: "100%",
                height: "100%"
            }}
        >
            {children}
        </div>
    );
}
const DndRecipient = ({ accept, onDragHover, onDrop, children }) => {
    const [{ isOver, canDrop }, drop] = useDrop({
        accept: accept,
        drop: (item, monitor) => {
            onDrop && onDrop(item, monitor);
        },
        hover: (item, monitor) => {
            onDragHover && onDragHover(item, monitor);
        },
        collect: (monitor) => ({
            isOver: !!monitor.isOver(),
            canDrop: !!monitor.canDrop(),
        }),
        canDrop: (item, monitor) => true
    });

    return (
        <div style={{
            position: 'relative', // 关键：设置相对定位容器
            width: "100%",
            height: "100%",
        }}>
            <Box
                ref={drop}
                sx={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    width: '100%',
                    height: '100%',
                    zIndex: 2, // 确保在按钮上方
                    pointerEvents: canDrop ? 'auto' : 'none', // 动态事件控制
                    // backgroundColor: isOver ? 'rgba(217, 0, 81, 0.1)' : 'transparent',
                }}
            />
            {children}
        </div>
    );
};


export { ItemTypes, DndSender, DndRecipient };