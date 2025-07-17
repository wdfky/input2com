import { Button, Grid, IconButton, Paper } from "@mui/material";
import { DndRecipient, ItemTypes } from "./dnd";
import CenterFocusWeakIcon from '@material-ui/icons/CenterFocusWeak';

const mouseCode2name = {
    1: "左键",
    2: "右键",
    4: "中间",
    8: "前进",
    16: "后退",
}

const MousePointer = ({ x, y }) => {
    return <IconButton
        style={{
            position: "absolute",
            left: x,
            top: y,
            width: "16px",
            height: "16px",
            transform: "translate(-50%, -50%)",
            color: "#757575"
        }}
    >
        <CenterFocusWeakIcon />
    </IconButton>
}


const MouseContainer = ({ setMouseConfig, onMouseClick }) => {

    return <div style={{
        height: "261px",
        position: "relative",
        padding: "24px",
        display: "inline-flex",
    }}>
        <img src="./mouse.png"
            style={{
                maxHeight: "100%",
                maxWidth: "120px",
                objectFit: "contain",
                width: "auto",
                height: "auto"
            }}
        />
        <Grid
            container
            direction="column"
            justifyContent="space-between"
            alignItems="flex-start"
            spacing="4px"
            sx={{
                width: "auto",
                height: "100%",
                marginLeft: "70px"
            }}
        >
            {[1 << 0, 1 << 1, 1 << 2, 1 << 3, 1 << 4].map((mouseCode) => (
                <Grid key={mouseCode} item xs={2} >
                    <DndRecipient accept={ItemTypes.CARD} onDragHover={() => { }} onDrop={(item, _) => { setMouseConfig(mouseCode, item["key"]) }} >
                        <Button
                            sx={{
                                width: 100,
                                height: "100%",
                                backgroundColor: "button.mouse.main",
                                "&:hover": {
                                    backgroundColor: "button.mouse.hover",
                                },
                                color: "button.mouse.text",
                                fontSize: "1rem",
                            }}

                            onClick={(e) => onMouseClick(mouseCode, e.clientX, e.clientY)}

                        >
                            {mouseCode2name[mouseCode]}
                        </Button>
                    </DndRecipient>
                </Grid>
            ))}
            <Grid key={"reset"} item xs={2} >
                <Button sx={{
                    width: 100,
                    height: "100%",
                    backgroundColor: "button.mouse.main",
                    "&:hover": {
                        backgroundColor: "button.mouse.hover",
                    },
                    color: "button.mouse.text",
                    fontSize: "1rem",

                }} onClick={() => {
                    fetch(`/api/set/mouse?key=${"CLEAR_ALL"}&value=NONE` )
                    fetch(`/api/set/keyboard?key=${"CLEAR_ALL"}&value=NONE` )
                }} >
                    清除所有
                </Button>
            </Grid>
        </Grid>
        <svg
            height={260}
            width={206}
            style={{
                position: "absolute",
                left: "0px",
                top: "0px",
            }}
        >
            <polyline points="52,53  100,30    206,30" fill="none" stroke="#757575" strokeWidth="1" />
            <polyline points="103,53  130,70    206,70" fill="none" stroke="#757575" strokeWidth="1" />
            <polyline points="78,70  130,105    206,105" fill="none" stroke="#757575" strokeWidth="1" />
            <polyline points="28,112  130,142    206,142" fill="none" stroke="#757575" strokeWidth="1" />
            <polyline points="28,148  130,178    206,178" fill="none" stroke="#757575" strokeWidth="1" />
        </svg>

        <MousePointer x={52} y={53} />
        <MousePointer x={103} y={53} />
        <MousePointer x={78} y={70} />
        <MousePointer x={28} y={112} />
        <MousePointer x={28} y={148} />

    </div>
}

export default MouseContainer;