import React from 'react'
import {Pagination} from 'antd'

// @ts-ignore
const CommonPagination = props => {
    return (
        <Pagination showQuickJumper={true}
                    showSizeChanger={true}
                    onChange={props.onChange}
                    pageSizeOptions={props.pageSizeOptions ? props.pageSizeOptions : ['10', '20', '50', '100', '200', '500']}
                    onShowSizeChange={props.onShowSizeChange}
                    showTotal={(total) => `共${total}条`}
                    {...props}
        />
    )
}
export default CommonPagination